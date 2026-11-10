import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const columns = [el.querySelector('path[d="M5 3v14"]'), el.querySelector('path[d="M19 3v18"]')];
  animate(columns, { scaleY: [0, 1.25, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(τmin / 2) });
  const card = el.querySelector('path[d="M12 3v8"]');
  animate(card, { translateX: [-7, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
