import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lights = ["M6 8h.01", "M10 8h.01", "M14 8h.01"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(lights, { scale: [1, 1.75, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
