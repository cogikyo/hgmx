import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const box = el.querySelectorAll('path[d^="M11 21.73"], path[d^="M12 22"], polyline');
  animate(box, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const tape = el.querySelector('path[d^="m7.5"]');
  animate(tape, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
