import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const pages = el.querySelector('path[d^="M20.001"]');
  animate(pages, { scaleX: [0.2, 1.1, 1] }, { duration: τmax });

  const spine = el.querySelector('path[d^="M12 5"]');
  animate(spine, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const lines = ["M6 9h2", "M6 13h2", "M16 9h2", "M16 13h2"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(lines, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
