import { animate, τmax } from "../motion.js";

export default function (el) {
  const top = el.querySelector('path[d="M4 5h16"]');
  animate(top, { translateY: [0, 4, -1, 0] }, { duration: τmax });

  const bottom = el.querySelector('path[d="M4 19h16"]');
  animate(bottom, { translateY: [0, -4, 1, 0] }, { duration: τmax });
}
