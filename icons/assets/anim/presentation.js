import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const rod = el.querySelector('path[d="M2 3h20"]');
  animate(rod, { scaleX: [0, 1] }, { duration: τmid });

  const screen = el.querySelector('path[d^="M21 3"]');
  animate(screen, { scaleY: [0, 1.15, 1] }, { duration: τmid, delay: τmin });

  const stand = el.querySelector('path[d^="m7 21"]');
  animate(stand, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
