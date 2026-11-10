import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const pages = el.querySelector('path[d^="M22 6"]');
  animate(pages, { scaleX: [0.2, 1.1, 1] }, { duration: τmax });

  const spine = el.querySelector('path[d^="M12 5"]');
  animate(spine, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const tick = el.querySelector('path[d^="m16 12"]');
  animate(tick, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmin });
  animate(tick, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
