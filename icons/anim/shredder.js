import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const strips = ["M6 20", "M10 22", "M14 19", "M18 20"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(strips, { translateY: [-3, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
