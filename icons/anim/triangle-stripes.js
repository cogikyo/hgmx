import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const sign = el.querySelector('path[d^="M13.75"]');
  animate(sign, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const stripes = ["M7.5 9h9", "M5.5 13h13", "M3 17h18"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(stripes, { scaleX: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
