import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const cup = el.querySelectorAll('ellipse, path[d^="M2 13v6"], path[d="M8 16v1"], path[d="M12 16v2"]');
  animate(cup, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const flame = el.querySelector('path[d^="M12 2s"]');
  animate(flame, { scaleY: [1, 0.8, 1.25, 1] }, { duration: τmid, delay: τmin });
}
