import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector('path[d^="M19 8"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const bar = el.querySelector('path[d="M12 2v20"]');
  animate(bar, { scaleY: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
