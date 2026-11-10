import { animate, stagger, τmin } from "../motion.js";

export default function (el) {
  const arrows = el.querySelectorAll("g");
  animate(arrows, { translateY: [6, -2, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12) });

  const line = el.querySelector('path[d="M4 21h16"]');
  animate(line, { translateY: [0, 1, 0] }, { duration: τmin });
}
