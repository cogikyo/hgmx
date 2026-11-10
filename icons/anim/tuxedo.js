import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bow = el.querySelector('path[d^="M10 3v2"]');
  animate(bow, { rotate: [0, -20, 10, 0] }, { duration: τmid });

  const buttons = el.querySelectorAll('path[d="M12 9h.01"], path[d="M12 13h.01"]');
  animate(buttons, { scale: [1, 1.25, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
