import { animate, stagger, τmin } from "../motion.js";

export default function (el) {
  const bulb = el.querySelector("circle");
  animate(bulb, { scale: [1, 0.85, 1] }, { duration: τmin });

  const mist = ["M5 4.5", "M2 2", "M2 7"].map((d) => el.querySelector(`path[d="${d}h.01"]`));
  animate(mist, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
