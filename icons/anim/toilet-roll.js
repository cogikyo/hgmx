import { animate } from "../motion.js";

export default function (el) {
  const perforation = el.querySelector('path[d="M10 14v2"]');
  animate(perforation, { translateY: [0, 1.5, 0] }, { duration: 0.3 });
}
