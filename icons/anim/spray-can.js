import { animate, stagger, τmin } from "../motion.js";

export default function (el) {
  const nozzle = el.querySelector("rect");
  animate(nozzle, { translateY: [0, 1, 0] }, { duration: τmin });

  const mist = ["M11 7", "M7 5", "M7 9", "M3 3", "M3 7", "M3 11"].map((d) => el.querySelector(`path[d="${d}h.01"]`));
  animate(mist, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03, { startDelay: τmin }) });
}
