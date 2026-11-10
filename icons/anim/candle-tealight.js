import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const cup = el.querySelectorAll('ellipse, path[d^="M2 13v6"]');
  animate(cup, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const wick = el.querySelector('path[d="M12 12V7"]');
  animate(wick, { scaleY: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
