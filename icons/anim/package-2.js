import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const box = el.querySelectorAll('path[d^="M16.76"], path[d^="M3.054"]');
  animate(box, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const tape = el.querySelector('path[d="M12 3v6"]');
  animate(tape, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
