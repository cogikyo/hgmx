import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const board = el.querySelector('path[d^="M16 4"]');
  animate(board, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const letter = el.querySelectorAll('path[d^="M9 12"], path[d="M12 11v6"]');
  animate(letter, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const foot = el.querySelector('path[d="M11 17h2"]');
  animate(foot, { scaleX: [0, 1.25, 1] }, { duration: τmin, delay: τmid });

  const clip = el.querySelector("rect");
  animate(clip, { translateY: [-3, 1, 0] }, { duration: τmin, delay: τmid });
}
