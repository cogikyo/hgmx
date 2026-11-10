import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const folder = el.querySelector('path[d^="M4 20"]');
  animate(folder, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const stem = el.querySelector('path[d="M12 15v5"]');
  animate(stem, { scaleY: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const node = el.querySelector("circle");
  animate(node, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
