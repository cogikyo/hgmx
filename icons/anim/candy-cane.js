import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const cane = el.querySelector('path[d^="M5.7"]');
  animate(cane, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const stripes = ["m4.874", "m7.906", "m10.8", "M17.75"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(stripes, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03, { startDelay: τmin }) });
}
