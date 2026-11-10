import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M6 22"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const cells = ["M8 13h2", "M14 13h2", "M8 17h2", "M14 17h2"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(cells, { scaleX: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
