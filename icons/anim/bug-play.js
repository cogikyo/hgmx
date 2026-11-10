import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bug = el.querySelector("g");
  animate(bug, { rotate: [0, -12, 10, -6, 0] }, { duration: τmid });

  const play = el.querySelector('path[d^="M14 15.003"]');
  animate(play, { translateX: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
