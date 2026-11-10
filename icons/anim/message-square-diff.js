import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bubble = el.querySelector('path[d^="M22 17"]');
  animate(bubble, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const plus = el.querySelectorAll('path[d="M10 9h4"], path[d="M12 7v4"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });

  const minus = el.querySelector('path[d="M10 15h4"]');
  animate(minus, { scaleX: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
