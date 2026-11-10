import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const arc = el.querySelector('path[d^="M12 2a"]');
  animate(arc, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const play = el.querySelector('path[d^="M9 9.003"]');
  animate(play, { translateX: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
