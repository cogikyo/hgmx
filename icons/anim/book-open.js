import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const pages = el.querySelector('path[d^="M20.001"]');
  animate(pages, { scaleX: [0.2, 1.1, 1] }, { duration: τmax });

  const spine = el.querySelector('path[d^="M12 5"]');
  animate(spine, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
}
