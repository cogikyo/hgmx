import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const arch = el.querySelector('path[d^="M18 20"]');
  animate(arch, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const bar = el.querySelector('path[d="M12 4v16"]');
  animate(bar, { scaleY: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
